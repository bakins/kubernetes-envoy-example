<?php
require_once __DIR__ . '/vendor/autoload.php';
require_once __DIR__ . '/GPBMetadata/Proto/Item/Item.php';
require_once __DIR__ . '/Item/Item.php';
require_once __DIR__ . '/Item/GetItemRequest.php';
require_once __DIR__ . '/Item/ListItemsRequest.php';
require_once __DIR__ . '/Item/ListItemsResponse.php';

// stub data
$items = array(
    "6ab9e0c2-e7be-4120-a3e9-62c39b7dbfd7" => array (
        "name" => "first item",
        "description" => "first item for sell"
    ),
    "4415fede-7462-4f12-b87f-ede596ec6ee2" => array (
        "name" => "second item",
        "description" => "another item for sell"
    ),
    "5d210689-cca7-4e81-8437-05b20f658ad0" => array (
        "name" => "item three",
        "description" => "Yet another item"
    ),
    "dff79aa1-6b13-4aeb-8dca-22a45322a293" => array (
        "name" => "item four",
        "description" => "something else"
    ),
    "6962f4ff-b752-4103-b90c-1f9bcec30913" => array (
        "name" => "fifth item",
        "description" => "last item for sell"
    )
);

$body = file_get_contents('php://input');

$request = new \Item\GetItemRequest();
$request->mergeFromString($body);

$id = $request->getId();
$item = $items[$id];

$response = new \Item\Item();
$response->setId($id);
$response->setName($item["name"]);
$response->setDescription($item["description"]);

header('Content-Type: application/grpc');
print($response->serializeToString());
