<?php

namespace App\Http\V1\Responses;

use Illuminate\Http\JsonResponse;
use Illuminate\Contracts\Support\Responsable;

class NoContent implements Responsable
{
    public static function responseCode(): int
    {
        return 204;
    }

    public function toResponse($request)
    {
        return new JsonResponse(
            "",
            self::responseCode(),
        );
    }

    public static function new(): self
    {
        return new self;
    }
}