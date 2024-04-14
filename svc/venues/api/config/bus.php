<?php


use App\Buses\DefaultBus;
use App\Commands\Middleware\ValidateCommand;

return [
    'buses' => [
        'default' => [
            'class' => DefaultBus::class,

            // So we could override to bind by a specific interface
            // 'bindAs' => DefaultBus::class,

            'preMiddleware' => [
                ValidateCommand::class,
            ],
            'postMiddleware' => [],
        ],
    ],
];