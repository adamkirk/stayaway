<?php

namespace App\Tests\Unit\Entities;

use Tests\UnitTestCase;
use App\ValueObjects\Uuid;
use App\Entities\Organisation;
use App\Exceptions\InvalidPropertyException;
use PHPUnit\Framework\Attributes\DataProvider;
use App\ValueObjects\Organisation as VO;

class OrganisationTest extends UnitTestCase
{
    protected function mockUuid(): Uuid
    {
        return $this->getMockBuilder(Uuid::class)
            ->disableOriginalConstructor()
            ->getMock();
    }

    public function test_successful_construction(): void
    {
        $uuid = $this->mockUuid();
        $name = $this->createMock(VO\Name::class);
        $slug = $this->createMock(VO\Slug::class);
        $org = Organisation::new($uuid, $name, $slug);

        $this->assertInstanceOf(Organisation::class, $org);
    }


    public function test_slug_defaults(): void
    {
        $uuid = $this->mockUuid();
        $name = $this->createMock(VO\Name::class);
        $name->expects($this->exactly(1))->method('value')->willReturn("A Valid Name !_special_chars");
        $org = Organisation::new($uuid, $name, null);

        $this->assertInstanceOf(Organisation::class, $org);
        $this->assertInstanceOf(VO\Slug::class, $org->slug());
        $this->assertEquals("a-valid-name-special-chars", $org->slug()->value());
    }


    public function test_getters_and_setters(): void
    {
        $uuid = $this->mockUuid();
        $name = $this->createMock(VO\Name::class);
        $slug = $this->createMock(VO\Slug::class);
        $org = Organisation::new($uuid, $name, $slug);

        $this->assertInstanceOf(Organisation::class, $org);
        $this->assertSame($uuid, $org->id());
        $this->assertSame($name, $org->name());
        $this->assertSame($slug, $org->slug());

        $newName = $this->createMock(VO\Name::class);
        $org->setName($newName);
        $this->assertSame($newName, $org->name());

        $newSlug = $this->createMock(VO\Slug::class);
        $org->setSlug($newSlug);
        $this->assertSame($newSlug, $org->slug());
    }
}
