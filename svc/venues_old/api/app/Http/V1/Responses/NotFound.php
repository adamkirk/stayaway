<?php

namespace App\Http\V1\Responses;

use App\Contracts\Arrayable;
use App\Api\Responses\ConvertsToJsonResponse;
use App\Api\Translation\TranslatesFieldNames;
use Illuminate\Contracts\Support\Responsable;

class NotFound implements Arrayable, Responsable
{
    use TranslatesFieldNames;
    use ConvertsToJsonResponse;

    public function __construct(
        public readonly string $message,
        public readonly string $code,
    ) {}

    public static function responseCode(): int
    {
        return 404;
    }

    public static function default(): self
    {
        return new self("resource not found", 404);
    }
}