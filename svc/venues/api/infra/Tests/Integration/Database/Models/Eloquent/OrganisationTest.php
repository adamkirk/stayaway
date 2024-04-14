<?php

namespace Infra\Tests\Integration\Database\Models\Eloquent;

use App\ValueObjects\Uuid;
use InvalidArgumentException;
use App\Entities\Organisation;
use Tests\IntegrationTestCase;
use Infra\Database\Models\Eloquent\Organisation as ElOrganisation;
use App\ValueObjects\Organisation as VO;

class OrganisationTest extends IntegrationTestCase
{
    public function test_model_is_inserted(): void
    {
        $id = Uuid::new()->toString();
        $model = new ElOrganisation;
        $model->id = $id;
        $model->name = "blah";
        $model->slug = "blah";
        $model->save();

        $this->assertDatabaseHas('organisations', [
            'id' => $id,
        ]);
    }

    public function test_model_can_be_converted_to_entity(): void
    {
        $id = Uuid::new()->toString();
        $model = new ElOrganisation;
        $model->id = $id;
        $model->name = "Valid Name";
        $model->slug = "valid-name";
        
        $entity = $model->toEntity();
        $this->assertInstanceOf(Organisation::class, $entity);
        $this->assertEquals($model->id, $entity->id()->toString());
        $this->assertEquals($model->name, $entity->name()->value());
        $this->assertEquals($model->slug, $entity->slug()->value());
    }

    public function test_it_can_be_converted_from_entity()
    {
        $id = Uuid::new();

        $nameValue = "A Valid Name";
        $name = $this->createMock(VO\Name::class);
        $name->expects($this->exactly(1))->method('value')->willReturn($nameValue);

        $slugValue = "a-valid-slug";
        $slug = $this->createMock(VO\Slug::class);
        $slug->expects($this->exactly(1))->method('value')->willReturn($slugValue);

        $entity = Organisation::new($id, $name, $slug);

        $model = ElOrganisation::fromEntity($entity);

        $this->assertEquals($model->id, $id->toString());
        $this->assertEquals($nameValue, $model->name);
        $this->assertEquals($slugValue, $model->slug);
    }

    public function test_it_can_be_updated_from_entity()
    {
        $id = Uuid::new();
        $model = new ElOrganisation;
        $model->id = $id->toString();
        $model->name = "blah";
        $model->slug = "blah";
        $model->save();

        $this->assertDatabaseHas('organisations', [
            'id' => $id,
        ]);

        $nameValue = "A Valid Name";
        $name = $this->createMock(VO\Name::class);
        $name->expects($this->exactly(1))->method('value')->willReturn($nameValue);

        $slugValue = "a-valid-slug";
        $slug = $this->createMock(VO\Slug::class);
        $slug->expects($this->exactly(1))->method('value')->willReturn($slugValue);
        $entity = Organisation::new($id, $name, $slug);

        $model->updateFromEntity($entity);

        $this->assertEquals($model->name, $nameValue);
        $this->assertEquals($model->slug, $slugValue);
    }

    public function test_it_cannot_be_updated_from_a_different_entity()
    {
        $id = Uuid::new()->toString();
        $model = new ElOrganisation;
        $model->id = $id;
        $model->name = "blah";
        $model->slug = "blah";
        $model->save();

        $this->assertDatabaseHas('organisations', [
            'id' => $id,
        ]);

        $id = Uuid::new();
        $name = $this->createMock(VO\Name::class);
        $name->expects($this->exactly(0))->method('value');

        $slug = $this->createMock(VO\Slug::class);
        $slug->expects($this->exactly(0))->method('value');
        $entity = Organisation::new($id, $name, $slug);

        $this->expectExceptionObject(new InvalidArgumentException("Cannot update model from an entity that has a different id!"));
        $model->updateFromEntity($entity);
    }
}
