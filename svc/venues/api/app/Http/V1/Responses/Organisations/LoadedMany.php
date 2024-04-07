<?php

namespace App\Http\V1\Responses\Organisations;

use App\Contracts\Arrayable;
use App\DTO\OrganisationPage;
use App\Http\V1\Schemas\ListMeta;
use App\Http\V1\Schemas\Organisation;
use App\Collections\OrganisationCollection;
use App\Api\Responses\ConvertsToJsonResponse;
use Illuminate\Contracts\Support\Responsable;
use App\Entities\Organisation as EOrganisation;

class LoadedMany implements Arrayable, Responsable
{
    use ConvertsToJsonResponse;

    protected function __construct(
        public readonly array $data,
        public readonly ListMeta $meta
    ) {}

    public static function responseCode(): int
    {
        return 200;
    }

    public static function fromPage(OrganisationPage $page): self
    {
        return new self(
            data: array_map(fn (EOrganisation $org) => Organisation::fromEntity($org), $page->organisations->all()),
            meta: ListMeta::fromDTOComponents($page->pagination),
        );
    }
}