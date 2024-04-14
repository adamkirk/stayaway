<?php

namespace App\Validation;

use App\Collections\ValidationErrorCollection;
use Symfony\Component\Validator\Validator\ValidatorInterface;

interface Validatable
{
    public function validate(ValidatorInterface $validator): ?ValidationErrorCollection;
}