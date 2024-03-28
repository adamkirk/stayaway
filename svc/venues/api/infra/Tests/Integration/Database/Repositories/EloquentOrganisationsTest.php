<?php

namespace Infra\Tests\Integration\Database\Repositories;

use App\ValueObjects\Uuid;
use InvalidArgumentException;
use App\Entities\Organisation;
use Tests\IntegrationTestCase;
use App\Repositories\SaveResult;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Infra\Database\Repositories\EloquentOrganisations;
use Infra\Database\Models\Eloquent\Organisation as ElOrganisation;

class EloquentOrganisationsTest extends IntegrationTestCase
{
    use RefreshDatabase;

    protected function getRepo(): EloquentOrganisations
    {
        return $this->app->make(EloquentOrganisations::class);
    }

    public function test_model_is_created()
    {
        $repo = $this->getRepo();
        $id = Uuid::new();

        $entity = Organisation::new($id, "Valid Name", "valid-slug");
        $result = $repo->save($entity);

        $this->assertEquals(SaveResult::Created, $result);
        $this->assertDatabaseHas('organisations', [
            'id' => $id->toString(),
        ]);
    }

    public function test_model_is_updated()
    {
        $repo = $this->getRepo();
        $id = Uuid::new();

        $entity = Organisation::new($id, "Valid Name", "valid-slug");
        $result = $repo->save($entity);

        $this->assertEquals(SaveResult::Created, $result);
        $this->assertDatabaseHas('organisations', [
            'id' => $id->toString(),
            'name' => 'Valid Name', 
        ]);

        $entity->setName("New Name");

        $result = $repo->save($entity);

        $this->assertEquals(SaveResult::Updated, $result);
        $this->assertDatabaseHas('organisations', [
            'id' => $id->toString(),
            'name' => 'New Name', 
        ]);
    }

    public function test_finds_model_by_id()
    {
        $repo = $this->getRepo();
        $id = Uuid::new();

        $entity = Organisation::new($id, "Valid Name", "valid-slug");
        $result = $repo->save($entity);

        $this->assertEquals(SaveResult::Created, $result);
        $this->assertDatabaseHas('organisations', [
            'id' => $id->toString(),
        ]);

        $entity = $repo->byId($id);

        $this->assertInstanceOf(Organisation::class, $entity);
        $this->assertEquals($id->toString(), $entity->id()->toString());
    }
}
