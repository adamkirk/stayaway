<?php

use App\Buses\DefaultBus;

return [
    'buses' => [
        'default' => [
            'class' => DefaultBus::class,

            // So we could override to bind by a specific interface
            // 'bindAs' => DefaultBus::class,

            'preMiddleware' => [],
            'postMiddleware' => [],
        ],
    ],
];