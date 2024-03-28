<?php

namespace App\Providers;

use App\Repositories\Organisations;
use Illuminate\Support\ServiceProvider;
use Symfony\Component\Validator\Validation;
use Infra\Database\Repositories\EloquentOrganisations;
use App\Http\Controllers\V1\Requests\CreateOrganisation;
use Symfony\Component\Validator\Validator\ValidatorInterface;

class AppServiceProvider extends ServiceProvider
{
    /**
     * Register any application services.
     */
    public function register(): void
    {
        $this->app->bind(ValidatorInterface::class, function($app) {
            return Validation::createValidatorBuilder()
                ->enableAttributeMapping()
                ->getValidator();
        });

        $this->app->singleton(Organisations::class, EloquentOrganisations::class);
    }

    /**
     * Bootstrap any application services.
     */
    public function boot(): void
    {
        //
    }
}
