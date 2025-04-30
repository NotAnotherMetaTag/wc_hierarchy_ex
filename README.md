# Go + React Engineering Challenge

Goal: Use requirements posited by Wider Circle to build a simple full-stack application comprised of a Golang back end and React front end.

## The App
The application will be stored together in this monorepo for ease of sharing and downloading. The Golang back end will be in the `server` directory. The front end will be stored in the `ui` directory.

### Back End
The back end is using Golang without any additional dependencies. It consumes mock data and does not store to any external database (data is only stored in memory). To launch the back end, run these commands:

```
cd server

go run .
```

You may need to allow access to your network.

###  Front End
The front end is using React.js with Vite, with testing via Vitest. To launch the front end, run these commands:

```
cd ui

npm install

npm run dev
```

Open the front end in your browser at `localhost:5173` (the Vite default).

I would recommend serving the back end prior to opening the front end, but it's built to be resilient.


### Testing
Executing Front End Unit Tests:
```
cd ui

npm test
```

Executing Back End Unit Tests:
```
cd server

go test
```

For some degree of convenience, I also hooked these up to my GitHub CI, so reports from these unit tests are readily available even if you did not feel like executing them yourself.

What the back end tests look for:
- Ensure Last Name is parsed correctly from first name. This includes people with two first names, special characters, and two last names.
- When a manager has multiple reports, those people are sorted alphabetically by last name.
- That the hierarchy is built correctly based on managerID mapping.
- Integration mock that parses the data, builds the hierarchy, and ultimately returns the sorted organization structure.
- That we do not support circular references, where people can be managers of one another.
- We check that we don't have duplicated IDs or invalid manager IDs.

**A Notable Gap In Testing**: Is when a person has a multi-word last name or a suffix, such as `Dick Van Dyke`, `Maria De La Cruz`, and `John Paul IV`. I opted for the absolute most simple method of targeting the last name based on where the last space is, to allow for multiple first names. It readily handles names with special characters, such as `O'Conner`, or hyphenated last names like `Olson-Williams`. 
The ways I considered implementing this to reduce edge case errors:
1. Instead look for the last, or second to last space for the last name. This fails my original test case of "Maria De La Cruz", but would have handled more of these edge cases correctly. On the flipside, this naive assertion would have then failed someone named something like `Anna Maria Perez`, with two first names. It could also cause confounding issues if anyone's name was provided with a middle name in addition to their first and last name, as it would be assumed part of the last name.
2. Create a list of possible edge case naming schemes, such as `"von", "Van", "De La"`, and more. This felt like trying to battle the world when our data set didn't even originally include any of these edge cases.

Ultimately, I decided to pick the absolute simplest option that solved my existing data set's needs. I would greatly prefer to in a production environment have this data already separated into `firstName`, `lastName`, `middleName`, `prefix`, and `suffix`. But if I had this data in production I think the labor intensive approach of the edge case naming schemes could be valuable if we truly couldn't split the data. And it would be labor intensive.

What the front end tests look for:
- The error/server not running state is handled gracefully.
- The retry button works as expected, appearing only when there is an error.
- The loading screen, when necessary, does show up.
- Shows the list when the GET call is successful.

To note, these tests all rely on "Michael Chen" as CEO, but we do hydrate this as mock data.

### Design Considerations:

1. Separation of Concerns
    - The back end handles the API call,  and error handling related to the API call. If the endpoint had provided data the front end did not need, I would have also transformed the received data here prior to sending it to the front end.
    - I did organize the hierarchy arrangement and data validation within the backend to keep the front end as simple as possible.  
    - The front end makes the request to the server to complete the API call, and displays the data it receives. I chose to have the backend sort the data in advance to keep the front end as simple as possible. Given I call the back end endpoint `/employees`, perhaps we could argue it makes sense to pass a parameter of `?sorted=true` to receive the hierarchal list (versus flat list), or move the hierarchy building to the front end. I think any of those options could be a valid way to tackle this requirement and in the context of a pre-existing app would defer to the existing design decisions.
2. Testability
    - The obvious ethos of unit testing is ensuring your functions "do one thing, and do it well". Given we are only doing one API call and displaying that data, there isn't a whole lot to test. But you should try to test things as you make them - it's way less work to do it now than to have to backtrack to add it later.
    - I considered whether it was worth adding some kind of automated testing suite like Selenium or Cypress for an end-to-end integration test, but decided to both not add a third "app" to run, as well as felt confident the unit tests were sufficient for the degree of complexity we have. 
    - The server does contain one integration test with a mocked payload and assertions against that payload, which I felt was good enough.
3. Frameworks (and lack of)
    - As I mentioned, I used Vite for creating the React App. Although I kind of liked create-react-app, it doesn't make sense to build something new with it. There are heavier frameworks that provide more built-in functionality that we could use if we had more requirements than our one API call. For WC's use-case, Expo could be especially valuable given the React native mobile application.
    - I didn't feel it necessary to enforce typescript in the front-end for the size of this application. I generally like typescript for enforcing more type safety, even if it does introduce more boilerplate code. Similarly, although Vite likes the SWC, in a project this small I felt no reason to include it.
    - There's literally nothing in the back end requirements for this challenge that Golang cannot do on its own. Just another case of keeping it simple. Something like Gin could be useful for managing a larger collection of APIs, cached sessions, and authentication.