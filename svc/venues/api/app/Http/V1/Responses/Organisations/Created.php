<?php

namespace App\Http\V1\Responses\Organisations;

use App\Contracts\Arrayable;
use App\Http\V1\Schemas\Organisation;
use App\Api\Responses\ConvertsToJsonResponse;
use Illuminate\Contracts\Support\Responsable;
use App\Entities\Organisation as EOrganisation;

class Created implements Arrayable, Responsable
{
    use ConvertsToJsonResponse;

    protected function __construct(
        public readonly Organisation $data,
    ) {}

    public static function responseCode(): int
    {
        return 201;
    }

    public static function fromEntity(EOrganisation $org): self
    {
        return new self(
            data: Organisation::fromEntity($org),
        );
    }
}