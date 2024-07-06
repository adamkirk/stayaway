<?php

namespace App\Tests\Unit\Collections;

use Tests\UnitTestCase;
use App\Entities\Organisation;
use App\Collections\OrganisationCollection;

class OrganisationCollectionTest extends UnitTestCase
{
    public function test_getting_all_items()
    {
        $orgs = [
            $this->createMock(Organisation::class),
            $this->createMock(Organisation::class),
            $this->createMock(Organisation::class),
        ];

        $coll = OrganisationCollection::fromArray($orgs);
        $this->assertEquals($orgs, $coll->all());

        $coll = OrganisationCollection::new(...$orgs);
        $this->assertEquals($orgs, $coll->all());

        $coll = OrganisationCollection::new();
        $coll->add($orgs[0]);
        $coll->add($orgs[1]);
        $coll->add($orgs[2]);
        $this->assertEquals($orgs, $coll->all());
    }

    /**
     * @covers App\Collections\ActsLikeArray
     */
    public function test_is_iterable()
    {
        $orgs = [
            $this->createMock(Organisation::class),
            $this->createMock(Organisation::class),
            $this->createMock(Organisation::class),
        ];

        $coll = OrganisationCollection::fromArray($orgs);

        $expectIndex = 0;
        foreach ($coll as $org) {
            $this->assertSame($orgs[$expectIndex], $org);
            $expectIndex++;
        }

        $this->assertEquals(3, $expectIndex, "did not iterate through collection");
    }


    /**
     * @covers App\Collections\ActsLikeArray
     */
    public function test_has_array_access()
    {
        $orgs = [
            $this->createMock(Organisation::class),
            $this->createMock(Organisation::class),
            $this->createMock(Organisation::class),
        ];

        $coll = OrganisationCollection::fromArray($orgs);

        $this->assertSame($orgs[0], $coll[0]);
        $this->assertSame($orgs[1], $coll[1]);
        $this->assertSame($orgs[2], $coll[2]);
    }
}
