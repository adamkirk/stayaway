<?php

namespace App\Api\Requests;

use App\Errors\ValidationErrorCollection;

interface Validatable
{
    public function validate(): ?ValidationErrorCollection;
}