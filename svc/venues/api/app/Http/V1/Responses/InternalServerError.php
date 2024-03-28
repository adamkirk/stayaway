<?php

namespace App\Http\V1\Responses;

use Ramsey\Uuid\Uuid;
use App\Contracts\Arrayable;
use App\Errors\ValidationErrorCollection;
use App\Api\Responses\ConvertsToJsonResponse;
use App\Api\Translation\TranslatesFieldNames;
use Illuminate\Contracts\Support\Responsable;

class InternalServerError implements Arrayable, Responsable
{
    use TranslatesFieldNames;
    use ConvertsToJsonResponse;

    public function __construct(
        public readonly string $requestId,
        public readonly string $message,
        public readonly string $code,
    ) {}

    public static function responseCode(): int
    {
        return 500;
    }

    public static function new(string $message = "Something went wrong!", string $code = "000"): self
    {
        // TODO: we need to be generating request ids somewhere
        // Then we can pull it from a header or something, just not implemented yet
        $requestId = Uuid::uuid4();

        return new self($requestId->toString(), $message, $code);
    }
}