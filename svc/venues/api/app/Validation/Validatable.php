<?php

namespace App\Validation;

use App\Collections\ValidationErrorCollection;
use Illuminate\Contracts\Support\Responsable;

interface Validatable
{
    public function validate(): ?ValidationErrorCollection;

    public function invalidResponse(ValidationErrorCollection $errors): Responsable;
}