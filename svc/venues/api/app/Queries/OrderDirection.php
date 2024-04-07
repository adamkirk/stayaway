<?php

namespace App\Queries;

enum OrderDirection: string
{
    const ASC = 'asc';
    const DESC = 'desc';
    const ALL = [
        self::ASC,
        self::DESC,
    ];

    case Ascending = self::ASC;
    case Descending = self::DESC;
}