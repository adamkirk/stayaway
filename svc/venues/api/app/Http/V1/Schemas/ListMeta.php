<?php

namespace App\Http\V1\Schemas;

use App\Contracts\Arrayable;
use App\Http\V1\Schemas\PaginationMeta;
use App\Api\Responses\ConvertsSelfToArray;
use App\Entities\Organisation as EOrganisation;
use App\DTO\PaginationMeta as PaginationMetaDTO;

class ListMeta implements Arrayable
{
    use ConvertsSelfToArray;

    protected function __construct(
        public readonly ?PaginationMeta $pagination,
    ) {}

    public static function fromDTOComponents(PaginationMetaDTO $pagination = null): self
    {
        return new self(
            pagination: $pagination !== null ? PaginationMeta::fromDTO($pagination) : null,
        );
    }
}