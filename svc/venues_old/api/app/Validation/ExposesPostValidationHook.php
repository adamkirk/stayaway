<?php

namespace App\Validation;

interface ExposesPostValidationHook
{
    public function postValidationHook(): void;
}