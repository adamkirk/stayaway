<?php

namespace App\Http\V1\Schemas;

use App\DTO\PaginationMeta as PaginationMetaDTO;
use App\Contracts\Arrayable;
use App\Api\Responses\ConvertsSelfToArray;

class PaginationMeta implements Arrayable
{
    use ConvertsSelfToArray;

    protected function __construct(
        public readonly int $page,
        public readonly int $pageSize,
        public readonly int $totalPages,
        public readonly int $totalResults,
    ) {}

    public static function fromDTO(PaginationMetaDTO $pagination): self
    {
        return new self(
            page: $pagination->page,
            pageSize: $pagination->pageSize,
            totalResults: $pagination->totalResults,
            totalPages: $pagination->totalPages,
        );
    }
}