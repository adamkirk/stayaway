<?php

namespace App\Handlers;

use App\Repositories\Organisations;
use App\Exceptions\NotFoundException;
use App\Commands\UpdateOrganisationCommand;

class UpdateOrganisationHandler
{
    public function __construct(
        protected Organisations $repo,
    ) {}

    /**
     * Handle the event.
     */
    public function handle(UpdateOrganisationCommand $cmd): void
    {
        $org = $this->repo->byId($cmd->id());

        if ($org === null) {
            throw new NotFoundException;
        }

        if (($name = $cmd->name()) !== null) {
            $org->setName($name);
        }

        if (($slug = $cmd->slug()) !== null) {
            $org->setSlug($slug);
        }

        $this->repo->save($org);
    }
}