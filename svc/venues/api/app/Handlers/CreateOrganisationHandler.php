<?php

namespace App\Handlers;

use App\Entities\Organisation;
use App\Repositories\Organisations;
use App\Commands\CreateOrganisationCommand;

class CreateOrganisationHandler
{
    public function __construct(
        protected Organisations $repo,
    ) {}

    /**
     * Handle the event.
     */
    public function handle(CreateOrganisationCommand $cmd): void
    {
        $org = Organisation::new($cmd->generatedId, $cmd->name, $cmd->slug);

        $this->repo->save($org);
    }
}