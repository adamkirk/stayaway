<?php

namespace App\Queries;

use Illuminate\Http\Request;
use App\Api\Translation\HttpField;
use App\Api\Translation\FieldPlacement;
use App\Validation\ExposesPostValidationHook;
use Symfony\Component\Validator\Constraints as Assert;

trait HasPaginationAndOrdering
{
    #[HttpField(in: FieldPlacement::Query)]
    #[Assert\GreaterThanOrEqual(value: 1)]
    #[Assert\LessThanOrEqual(value: 100)]
    public readonly int $page;

    #[HttpField(in: FieldPlacement::Query)]
    #[Assert\GreaterThanOrEqual(value: 1)]
    public readonly int $pageSize;

    #[HttpField(in: FieldPlacement::Query)]
    #[Assert\Regex('/^[A-Za-z0-9_]+$/')]
    #[Assert\NotBlank]
    public readonly ?string $orderBy;

    #[HttpField(name: "order_dir", in: FieldPlacement::Query)]
    #[Assert\Choice(choices: OrderDirection::ALL)]
    #[Assert\NotBlank]
    public readonly ?string $rawOrderDirection;

    public readonly OrderDirection $orderDirection;

    protected function populatePaginationAndOrdering(Request $request): void
    {
        $this->page = $request->get($this->translate('page'), self::getDefaultPage());
        $this->pageSize = $request->get($this->translate('pageSize'), self::getDefaultPageSize());

        $this->orderBy = $request->get($this->translate('orderBy'), self::getDefaultOrderField());
        $this->rawOrderDirection = $request->get($this->translate('rawOrderDirection'), self::getDefaultOrderDirection());
    }

    protected static function getDefaultPage(): int
    {
        return 1;
    }

    protected static function getDefaultPageSize(): int
    {
        return 20;
    }

    public function paginationAndOrderingPostHook(): void
    {
        $this->orderDirection = OrderDirection::from($this->rawOrderDirection);
    }

    abstract protected static function getDefaultOrderField(): string;

    abstract protected static function getDefaultOrderDirection(): string;
}