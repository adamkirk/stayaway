<?php

namespace App\DTO;

class PaginationMeta
{
    public function __construct(
        public readonly int $page,
        public readonly int $pageSize,
        public readonly int $totalPages,
        public readonly int $totalResults,
    ){}
}