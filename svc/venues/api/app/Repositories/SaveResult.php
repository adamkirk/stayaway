<?php

namespace App\Repositories;

enum SaveResult 
{
    case Created;
    case Updated;
    case Failed;
    case DidNothing;
}