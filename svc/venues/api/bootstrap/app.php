<?php

use Illuminate\Foundation\Application;
use Illuminate\Foundation\Configuration\Exceptions;
use Illuminate\Foundation\Configuration\Middleware;
use Illuminate\Foundation\Http\Middleware\ConvertEmptyStringsToNull;

return Application::configure(basePath: dirname(__DIR__))
    ->withRouting(
        web: __DIR__.'/../routes/web.php',
        api: __DIR__.'/../routes/api.php',
        commands: __DIR__.'/../routes/console.php',
        health: '/up',
    )
    ->withEvents(discover: [
        __DIR__.'/../app/Handlers',
    ])
    ->withMiddleware(function (Middleware $middleware) {
        // This only confuses things in request validation, null and an empty
        // string are not the same thing...
        $middleware->remove(ConvertEmptyStringsToNull::class);
    })
    ->withExceptions(function (Exceptions $exceptions) {
        //
    })->create();
