<?php

namespace App\Providers;

use App\Buses\Bus;
use League\Tactician\CommandBus;
use App\Repositories\Organisations;
use App\Buses\DefinedHandlerLocator;
use Symfony\Component\Finder\Finder;
use Illuminate\Support\ServiceProvider;
use Symfony\Component\Validator\Validation;
use Illuminate\Contracts\Foundation\Application;
use Infra\Database\Repositories\EloquentOrganisations;
use League\Tactician\Handler\CommandHandlerMiddleware;
use App\Http\Controllers\V1\Requests\CreateOrganisation;
use Symfony\Component\Validator\Validator\ValidatorInterface;
use League\Tactician\Handler\MethodNameInflector\HandleInflector;
use League\Tactician\Handler\CommandNameExtractor\ClassNameExtractor;
use League\Tactician\Handler\MethodNameInflector\HandleClassNameInflector;

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

        $this->setupCommandBuses();
    }

    protected function setupCommandBuses(): void
    {
        foreach (config('bus.buses') as $name => $config) {
            $class = $config['class'];

            $implements = class_implements($class);
            if (! in_array(Bus::class, $implements)) {
                throw new \Exception("Buses must implement:" . Bus::class);
            }

            $bindAs = $config['bindAs'] ?? $class;

            // i'd expect this function is only called once on app boot (octane and all that)
            // But that doesn't seem to the the case, this is constructed for every request it seems
            // Not sure if this a bug, some dodgy setup on my part, or q nuance i'm missing...
            $this->app->singleton($bindAs, function(Application $app) use($config, $class, $bindAs) {
                // See: https://tactician.thephpleague.com/tweaking-tactician/
                $commandHandlerMiddleware = new CommandHandlerMiddleware(
                    new ClassNameExtractor,
                    new DefinedHandlerLocator, 
                    new HandleInflector,
                );

                $pre = array_map(fn($middleware) => app($middleware), $config['preMiddleware']);
                $post = array_map(fn($middleware) => app($middleware), $config['postMiddleware']);

                // To stop recursive loop when setting up binding
                $bus = $bindAs != $class ? app($config['class']) : new $class;

                $bus->setCommandBus(new CommandBus(
                    $pre + [$commandHandlerMiddleware] + $post,
                ));

                return $bus;
            });
        }
    }

    /**
     * Bootstrap any application services.
     */
    public function boot(): void
    {
        //
    }
}
