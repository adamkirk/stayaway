<?php

namespace App\Tests\Unit\ValueObjects;

use Tests\UnitTestCase;
use App\ValueObjects\Uuid;
use InvalidArgumentException;
use Ramsey\Uuid\Rfc4122\Fields;
use Ramsey\Uuid\Uuid as RamseyUuid;

class UuidTest extends UnitTestCase
{
    public function test_that_it_builds_a_uuidv7()
    {
        $vo = Uuid::new();

        $uuid = RamseyUuid::fromString($vo->toString());

        $fields = new Fields($uuid->getFields()->getBytes());

        $this->assertEquals($fields->getVersion(), 7);
    }

    public function test_that_it_builds_a_uuidv7_from_string()
    {
        $vo = Uuid::fromString(RamseyUuid::uuid7()->toString());

        $uuid = RamseyUuid::fromString($vo->toString());

        $fields = new Fields($uuid->getFields()->getBytes());

        $this->assertEquals($fields->getVersion(), 7);
    }

    public function test_that_it_throws_if_using_non_v7_uuid()
    {
        $this->expectExceptionObject(new InvalidArgumentException("Must supply a uuid v7"));
        Uuid::fromString(RamseyUuid::uuid4()->toString());
    }
}
