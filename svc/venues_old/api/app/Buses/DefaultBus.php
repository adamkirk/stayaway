<?php

namespace App\Buses;

use App\Buses\Bus;
use League\Tactician\CommandBus;

class DefaultBus implements Bus
{
    protected CommandBus $tactician;

    public function __construct() {}

    public function setCommandBus(CommandBus $bus): void
    {
        $this->tactician = $bus;
    }

    public function handle($cmd): void
    {
        $this->tactician->handle($cmd);
    }

    public function queue($cmd): void
    {
        // queue the command
        // $this->tactician->handle($cmd);
    }
}