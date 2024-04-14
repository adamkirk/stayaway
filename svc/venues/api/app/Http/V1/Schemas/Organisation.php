<?php

namespace App\Http\V1\Schemas;

use App\Contracts\Arrayable;
use App\Api\Responses\ConvertsSelfToArray;
use App\Entities\Organisation as EOrganisation;

class Organisation implements Arrayable
{
    use ConvertsSelfToArray;

    protected function __construct(
        public readonly string $id,
        public readonly string $name,
        public readonly string $slug,
    ) {}

    public static function fromEntity(EOrganisation $org): self
    {
        return new self(
            id: $org->id()->toString(),
            name: $org->name()->value(),
            slug: $org->slug()->value(),
        );
    }


}