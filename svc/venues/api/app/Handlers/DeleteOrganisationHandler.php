<?php

namespace App\Handlers;

use App\Entities\Organisation;
use App\Repositories\DeleteResult;
use App\Repositories\Organisations;
use App\Exceptions\NotFoundException;
use App\Commands\CreateOrganisationCommand;
use App\Commands\DeleteOrganisationCommand;

class DeleteOrganisationHandler
{
    public function __construct(
        protected Organisations $repo,
    ) {}

    /**
     * Handle the event.
     */
    public function handle(DeleteOrganisationCommand $cmd): void
    {
        $result = $this->repo->delete($cmd->id);

        if ($result == DeleteResult::NotFound) {
            throw new NotFoundException;
        }
    }
}