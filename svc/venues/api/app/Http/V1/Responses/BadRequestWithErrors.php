<?php

namespace App\Http\V1\Responses;

use App\Contracts\Arrayable;
use App\Collections\ValidationErrorCollection;
use App\Api\Responses\ConvertsToJsonResponse;
use App\Api\Translation\TranslatesFieldNames;
use Illuminate\Contracts\Support\Responsable;

class BadRequestWithErrors implements Arrayable, Responsable
{
    use TranslatesFieldNames;
    use ConvertsToJsonResponse;

    public function __construct(
        public readonly array $errors,
    ) {}

    public static function responseCode(): int
    {
        return 400;
    }

    public static function new(ValidationErrorCollection $errorCollection): self
    {
        $errors = [];

        foreach ($errorCollection->all() as $value) {
            $errors[$value->name()][] = $value->message();
        }

        return new self($errors);
    }
}