<?php

namespace App\DTO;

use App\Collections\OrganisationCollection;

class OrganisationPage
{
    public function __construct(
        public readonly OrganisationCollection $organisations,
        public readonly PaginationMeta $pagination
    ){}
}