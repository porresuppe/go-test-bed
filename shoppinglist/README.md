Exercise - Shopping List

Living on a student budget can be tight sometimes. You want to build a service that allows you to create a smart shopping list with the things you need and the supermarket where you can find them for the cheapest price. So let's build a web service for that:

The application should allow you to:

    Add a new item to the list
    Remove a single item
    Remove all items
    Return the total price for all your items
    Return all items that you need to buy in a single supermarket

You should have HTTP endpoints for all above. For example you might have: /items that returns all items currently on the list if you send a GET request and allows you to add a new item when sending a POST request with a JSON object.

A JSON object for a shopping item might have the following structure:

{
    "name": "milk",
    "supermarket" : "netto",
    "price" : 10.5
}

Note: It's fine to keep all the data in memory for this exercise.