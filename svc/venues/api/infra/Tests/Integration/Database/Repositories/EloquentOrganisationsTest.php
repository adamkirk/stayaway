<?php

namespace Infra\Tests\Integration\Database\Repositories;

use App\ValueObjects\Uuid;
use InvalidArgumentException;
use App\Entities\Organisation;
use Tests\IntegrationTestCase;
use App\Queries\OrderDirection;
use App\Repositories\SaveResult;
use App\Repositories\DeleteResult;
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


    public function test_page()
    {
        $repo = $this->getRepo();

        $toCreate = 25;
        for ($i=1; $i <= $toCreate; $i++) {
            $id = Uuid::new();
            $result = $repo->save(Organisation::new($id, "Valid Name $i", "valid-slug-$i"));
            $this->assertEquals(SaveResult::Created, $result);
        }

        
        $this->assertDatabaseCount('organisations', $toCreate);

        $pageNumber = 1;
        $pageSize = 10;
        $orderField = 'name';
        $orderDir = OrderDirection::Ascending;

        $page = $repo->page($pageNumber, $pageSize, $orderField, $orderDir);

        $this->assertEquals($pageNumber, $page->pagination->page);
        $this->assertEquals($pageSize, $page->pagination->pageSize);
        $this->assertEquals($toCreate, $page->pagination->totalResults);
        $this->assertEquals(3, $page->pagination->totalPages);

        $slugs = array_map(fn(Organisation $org) => $org->slug(), $page->organisations->all());

        /*
        Natural sorting in MySQL is a pain, so this will do for now. As long as 
        it's ordering them by the correct field (however mysql wants to do that)
        */
        $this->assertEquals([
            'valid-slug-1',
            'valid-slug-10',
            'valid-slug-11',
            'valid-slug-12',
            'valid-slug-13',
            'valid-slug-14',
            'valid-slug-15',
            'valid-slug-16',
            'valid-slug-17',
            'valid-slug-18',
        ], $slugs);


        // Change page
        $pageNumber = 2;
        $page = $repo->page($pageNumber, $pageSize, $orderField, $orderDir);

        $this->assertEquals($pageNumber, $page->pagination->page);
        $this->assertEquals($pageSize, $page->pagination->pageSize);
        $this->assertEquals($toCreate, $page->pagination->totalResults);
        $this->assertEquals(3, $page->pagination->totalPages);

        $slugs = array_map(fn(Organisation $org) => $org->slug(), $page->organisations->all());

        $this->assertEquals([
            'valid-slug-19',
            'valid-slug-2',
            'valid-slug-20',
            'valid-slug-21',
            'valid-slug-22',
            'valid-slug-23',
            'valid-slug-24',
            'valid-slug-25',
            'valid-slug-3',
            'valid-slug-4',
        ], $slugs);

        // Final page
        $pageNumber = 3;
        $page = $repo->page($pageNumber, $pageSize, $orderField, $orderDir);

        $this->assertEquals($pageNumber, $page->pagination->page);
        $this->assertEquals($pageSize, $page->pagination->pageSize);
        $this->assertEquals($toCreate, $page->pagination->totalResults);
        $this->assertEquals(3, $page->pagination->totalPages);

        $slugs = array_map(fn(Organisation $org) => $org->slug(), $page->organisations->all());

        $this->assertEquals([
            'valid-slug-5',
            'valid-slug-6',
            'valid-slug-7',
            'valid-slug-8',
            'valid-slug-9',
        ], $slugs);

        // Page out of bounds, but won't error here
        $pageNumber = 4;
        $page = $repo->page($pageNumber, $pageSize, $orderField, $orderDir);

        $this->assertEquals($pageNumber, $page->pagination->page);
        $this->assertEquals($pageSize, $page->pagination->pageSize);
        $this->assertEquals($toCreate, $page->pagination->totalResults);
        $this->assertEquals(3, $page->pagination->totalPages);

        $this->assertTrue($page->organisations->isEmpty());

        // Change page size and order
        $pageNumber = 1;
        $pageSize = 3;
        $orderDir = OrderDirection::Descending;
        $page = $repo->page($pageNumber, $pageSize, $orderField, $orderDir);

        $this->assertEquals($pageNumber, $page->pagination->page);
        $this->assertEquals($pageSize, $page->pagination->pageSize);
        $this->assertEquals($toCreate, $page->pagination->totalResults);
        $this->assertEquals(9, $page->pagination->totalPages);

        $slugs = array_map(fn(Organisation $org) => $org->slug(), $page->organisations->all());

        $this->assertEquals([
            'valid-slug-9',
            'valid-slug-8',
            'valid-slug-7',
        ], $slugs);

        // Next page
        $pageNumber = 2;
        $pageSize = 3;
        $orderDir = OrderDirection::Descending;
        $page = $repo->page($pageNumber, $pageSize, $orderField, $orderDir);

        $this->assertEquals($pageNumber, $page->pagination->page);
        $this->assertEquals($pageSize, $page->pagination->pageSize);
        $this->assertEquals($toCreate, $page->pagination->totalResults);
        $this->assertEquals(9, $page->pagination->totalPages);

        $slugs = array_map(fn(Organisation $org) => $org->slug(), $page->organisations->all());

        $this->assertEquals([
            'valid-slug-6',
            'valid-slug-5',
            'valid-slug-4',
        ], $slugs);
    }

    public function test_delete()
    {
        $repo = $this->getRepo();
        $id = Uuid::new();

        $entity = Organisation::new($id, "Valid Name", "valid-slug");
        $result = $repo->save($entity);

        $this->assertEquals(SaveResult::Created, $result);
        $this->assertDatabaseHas('organisations', [
            'id' => $id->toString(),
        ]);

        $result = $repo->delete($id);

        $this->assertEquals(DeleteResult::Deleted, $result);
        $this->assertDatabaseMissing('organisations', [
            'id' => $id->toString(),
        ]);

        // Can't delete if it's not there
        $result = $repo->delete($id);
        $this->assertEquals(DeleteResult::NotFound, $result);
    }
}
