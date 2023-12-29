

from typing import List


# Example usage
class Entity:
    def __init__(self, name: str):
        self.name = name


# Create a generic Repository class
class Repository[T:Entity]:
    def __init__(self, entity: T):
        self.entities: List[entity] = []
        print(entity.__name__)

    def add_entity(self, entity: T):
        self.entities.append(entity)

    def get_all_entities(self) -> List[T]:
        return self.entities


# Example usage
class User(Entity):
    def __init__(self, name: str):
        super().__init__(name)
        self.name = name


class Car:
    def __init__(self, name: str):
        self.name = name


def a(a: Repository[Car]):
    pass


car_repository = Repository[Car](Car)
user_repository = Repository[User](User)
user_repository.add_entity(User("Alice"))
car_repository.add_entity(
    Car("Bob"))  # this give warning: Expected type 'Car' (matched generic type 'T'), got 'Car' instead
a(car_repository)
a(user_repository)  # this give warning: Expected type 'Repository[Car]', got 'Repository[User]' instead
users = user_repository.get_all_entities()
print(users)
