<?php

namespace App\Tests\Integration\Commands;

use App\Errors\ErrorType;
use App\ValueObjects\Uuid;
use Illuminate\Http\Request;
use Illuminate\Routing\Route;
use Tests\IntegrationTestCase;
use App\Errors\ValidationError;
use App\Commands\DeleteOrganisationCommand;
use PHPUnit\Framework\MockObject\MockObject;
use App\Collections\ValidationErrorCollection;
use Symfony\Component\Validator\Validator\ValidatorInterface;

class DeleteOrganisationCommandTest extends IntegrationTestCase
{
    protected function getValidator(): ValidatorInterface
    {
        return resolve(ValidatorInterface::class);
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

    public function test_fields_are_populated_if_set()
    {
        $subj = new DeleteOrganisationCommand($this->getValidator());

        $uuid = '018eb897-4323-76b3-9c55-483ab7f55f43';
        $route = $this->getMockRoute();
        $route->expects($this->exactly(1))->method('parameter')->with('organisation_id')->willReturn($uuid);
        $req = $this->getMockRequest($route);

        $subj->populate($req);
        // Sets the id property
        $subj->postValidationHook();

        $this->assertEquals($uuid, $subj->rawId);
        $this->assertInstanceOf(Uuid::class, $subj->id);
        $this->assertEquals($uuid, $subj->id->toString());
    }

    /**
     * @dataProvider invalidIdValues
     */
    public function test_id_must_be_uuidv7(mixed $id)
    {
        $subj = new DeleteOrganisationCommand($this->getValidator());

        $route = $this->getMockRoute();
        $route->expects($this->exactly(1))->method('parameter')->with('organisation_id')->willReturn($id);
        $req = $this->getMockRequest($route);

        $subj->populate($req);
        $errors = $subj->validate();
        $this->assertInstanceOf(ValidationErrorCollection::class, $errors);

        $this->assertEquals(
            new ValidationError(
                name: "organisation_id",
                errorType: ErrorType::ValueNotAllowed,
                message: "This is not a valid UUID v7.",
            ),
            $errors[0],
        );
    }

    public static function invalidIdValues(): array
    {
        return [
            'uuidv4' => ['12242d10-84c4-49ef-8979-ee40013992fc'],
            'random-string' => ['blah'],
            'int' => [12345],
        ];
    }
}
