<?php

namespace App\Api\Requests;

interface ExposesPostValidationHook
{
    public function postValidationHook(): void;
}