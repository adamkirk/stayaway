<?php

namespace App\Buses;

use League\Tactician\CommandBus;

interface Bus
{
    public function setCommandBus(CommandBus $bus): void;

    public function handle($cmd): void;
}