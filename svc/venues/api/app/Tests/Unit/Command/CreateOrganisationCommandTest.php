<?php

namespace App\Tests\Entities;

use Tests\UnitTestCase;
use App\ValueObjects\Uuid;
use Illuminate\Http\Request;
use Illuminate\Routing\Route;
use App\Commands\CreateOrganisationCommand;
use PHPUnit\Framework\MockObject\MockObject;
use Symfony\Component\Validator\Validator\ValidatorInterface;

class CreateOrganisationCommandTest extends UnitTestCase
{
    protected function getMockValidator(): MockObject
    {
        return $this->getMockBuilder(ValidatorInterface::class)
            ->disableOriginalConstructor()
            ->getMock();
    }

    protected function getMockRoute(): MockObject
    {
        return $this->getMockBuilder(Route::class)
            ->disableOriginalConstructor()
            ->getMock();
    }

    protected function getMockRequest(MockObject $route): MockObject
    {
        $req = $this->getMockBuilder(Request::class)
            ->disableOriginalConstructor()
            ->getMock();

        $req->expects($this->any())->method('route')->willReturn($route);

        return $req;
    }

    /**
     * The generatedId should be created each time on construct
     */
    public function test_fields_are_populated_if_set()
    {
        $subj = new CreateOrganisationCommand($this->getMockValidator());

        $route = $this->getMockRoute();
        $req = $this->getMockRequest($route);
        $req->expects($this->exactly(2))->method('get')->willReturnCallback(function($key, $default) {
            return match($key) {
                'name' => 'my name',
                'slug' => 'some-slug',
            };
        });

        $subj->populate($req);
        $this->assertInstanceOf(Uuid::class, $subj->generatedId);
        $this->assertEquals($subj->name, 'my name');
        $this->assertEquals($subj->slug, 'some-slug');
    }
}
