<?php

namespace App\Buses;

use App\Buses\DefinesHandler;
use League\Tactician\Handler\Locator\HandlerLocator;

class DefinedHandlerLocator implements HandlerLocator
{
    public function getHandlerForCommand($cmd)
    {
        if (!class_exists($cmd)) {
            throw new \Exception("$cmd does not exist!");
        }

        $implements = class_implements($cmd);

        if (! in_array(DefinesHandler::class, $implements)) {
            throw new \Exception("$cmd must implement " . DefinesHandler::class);
        }

        // Hate the app function, but only other way is to inject the container
        // Octane says this might be a problem? so meh...
        // https://laravel.com/docs/9.x/octane#container-injection
        return app($cmd::getHandler());
    }
}