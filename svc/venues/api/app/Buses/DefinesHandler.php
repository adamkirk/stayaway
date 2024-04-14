<?php

namespace App\Buses;

interface DefinesHandler
{
    public static function getHandler(): string;
}