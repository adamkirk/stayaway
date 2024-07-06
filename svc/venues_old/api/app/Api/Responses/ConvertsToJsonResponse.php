<?php

namespace App\Api\Responses;

use Illuminate\Http\JsonResponse;

trait ConvertsToJsonResponse
{
    use ConvertsSelfToArray;

    abstract public static function responseCode(): int;

    public function toResponse($request)
    {
        return new JsonResponse(
            $this->toArray(),
            self::responseCode(),
        );
    }
}