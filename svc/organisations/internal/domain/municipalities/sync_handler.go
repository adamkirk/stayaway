package municipalities

import (
	"sync"

	"github.com/adamkirk-stayaway/organisations/internal/domain/common"
	"github.com/gocarina/gocsv"
	"github.com/spf13/afero"
)

type SyncHandlerConfig interface {
	MunicipalitiesSyncBatchSize() int
	MunicipalitiesSyncMaxProcesses() int
	MunicipalitiesSyncCountries() []string
}

type SyncHandlerRepo interface {
	UpdateBatch(batch []Municipality) (BatchUpdateResult, error)
}

type SyncHandler struct {
	fs afero.Fs
	repo SyncHandlerRepo
	cfg SyncHandlerConfig
}

type MunicipalityBatch struct {
	Number int
	Rows []Municipality
}

type SyncCommand struct {
	SourceCsvPath string
}

type batchResult struct {
	SyncedRecords int
	Error error
}

type processorFunc func (rows MunicipalityBatch, wg *sync.WaitGroup, sem processorSemaphore, resChan chan<- batchResult) 

type processorSemaphore chan MunicipalityBatch

func (h *SyncHandler) shouldIncludeRow(row Municipality) bool {
	for _, c := range h.cfg.MunicipalitiesSyncCountries() {
		if c == row.Country {
			return true
		}
	}

	return false
}

func (h *SyncHandler) processBatches(path string, processor processorFunc) ([]batchResult, error) {

	f, err := h.fs.Open(path)

	if err != nil {
		return nil, err
	}
	defer f.Close()


	rows := []Municipality{}

	if err := gocsv.Unmarshal(f, &rows); err != nil {
		return nil, err
	}

	batchSize := h.cfg.MunicipalitiesSyncBatchSize()
	batches := [][]Municipality{}
	

	batch := []Municipality{}
	for _, row := range rows {
		if ! h.shouldIncludeRow(row) {
			continue
		}

		batch = append(batch, row)

		if len(batch) == batchSize {
			batches = append(batches, batch)
			batch = []Municipality{}
		}
	}

	// Catch any leftovers in the final partial batch
	if len(batch) > 0 {
		batches = append(batches, batch)
	}

	resultsChannel := make(chan batchResult, len(batches))
	var wg sync.WaitGroup
	sem := make(processorSemaphore, h.cfg.MunicipalitiesSyncMaxProcesses()) // Semaphore to limit to X concurrent processes

	for i, batch := range batches {
		wg.Add(1)
		go processor(MunicipalityBatch{
			Number: i,
			Rows: batch,
		}, &wg, sem, resultsChannel)
	}

	wg.Wait()
	close(resultsChannel)

	results := []batchResult{}
	
	for res := range resultsChannel {
		results = append(results, res)
	}

	return results, nil
}

func (h *SyncHandler) syncBatch(batch MunicipalityBatch, wg *sync.WaitGroup, sem processorSemaphore, resChan chan<- batchResult) {
	defer wg.Done()

	sem <- batch // Acquire semaphore

	defer func() { <-sem }() // Release semaphore

	updateResult, err := h.repo.UpdateBatch(batch.Rows)

	res := batchResult{
		SyncedRecords: updateResult.Created + updateResult.Updated,
		Error: err,
	}

	resChan <- res
}

func (h *SyncHandler) Handle(cmd SyncCommand) (SyncResult, error) {
	res, err := h.processBatches(cmd.SourceCsvPath, h.syncBatch)

	if err != nil {
		return SyncResult{}, err
	}

	totalProcessed := 0
	errs := []error{}

	for _, b := range res {
		totalProcessed += b.SyncedRecords
		if b.Error != nil {
			errs = append(errs, b.Error)
		}
	}

	if len(errs) > 0 {
		return SyncResult{}, common.ErrGroup{
			Errors: errs,
		}
	}

	return SyncResult{
		Processed: totalProcessed,
		Path: cmd.SourceCsvPath,
	}, nil
}

func NewSyncHandler(fs afero.Fs, cfg SyncHandlerConfig, repo SyncHandlerRepo) *SyncHandler {
	return &SyncHandler{
		fs: fs,
		repo: repo,
		cfg: cfg,
	}
}