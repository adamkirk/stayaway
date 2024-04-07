<?php

namespace App\Api\Requests;

use App\Errors\ValidationErrorCollection;
use Illuminate\Contracts\Support\Responsable;

interface Validatable
{
    public function validate(): ?ValidationErrorCollection;

    public function invalidResponse(ValidationErrorCollection $errors): Responsable;
}