
#include <iostream>
#include <string>



class  Book {
    private:
        std::string name;
        std::string author;

    
    public:
        int year;
        // default constructor takes no argument
        // uses initializer list to set values to the argument
        Book() : name("vinod"),author("yatra"), year(2007){}


        // copy constructor -- creates deepcopy of the object
        // Syntax :: Classname(const Classname& other);
        // Should take a constant reference to the object.
        Book(const Book& other): name(other.name),author(other.author), year(other.year) {}

        // copy assignment operator
        // We are defining the copying behaviour of the operation My_second_book = My_first_book.
        // "operator=" --> syntax for performing operator overloading of assignment operator "="
        Book& operator=(const Book& other){
            // Check for self assignment (eg: b = b), during self assignment returns nothing.
            // this => special pointer to the current object.
            if (this != &other){
                this -> name = other.name;
                this -> author = other.author;
                this -> year = other.year;
            }

            return *this;
        }

        // move constructor performs two tasks:
        // - steal: Transfer resources from other to this.
        // - nullify: set the resources from other object to a safe, empty state.

        // std::move() is a function in cpp which performs this function 
        // move semantics are built in to all standard data types which is 
        // invoked by the std::move() function.

        // for cutsom data types/classes, move constructor is called when a new object
        // is called from a temporary object (rvalue reference.)
        //      ClassName(ClassName&&, other);
        Book(Book&& other) : name(std::move(other.name)), author(std::move(other.author)), year(std::move(other.year)+1){};

        // move-assingment constructor

        Book& operator=(Book&& other){
            if (this != &other){
                this -> name = std::move(other.name);
                this -> author = std::move(other.author);
                this -> year = std::move(other.year) + 10;

            };
            return *this;
        }
                             


        // destructor 
        // do not specify any arguments as data types like
        // int and string are automatically handled by compiler
        ~Book(){};


        
};

int main(){

    Book my_book;
    

    // invoking copy constructor
    Book my_second_book = my_book;
    // Book my_third_book(my_book);

    // copy-assignment operator

    Book my_third_book;
    my_third_book = my_book;

    // invoking move constructor
    Book my_fourth_book(std::move(my_second_book));
    // Book my_fourth_book = std::move(my_second_book)

    Book my_fifth_book;
    my_fifth_book = std::move(my_book);



    std::cout << my_third_book.year << "\n";
    std::cout << my_fourth_book.year << "\n";
    std::cout << my_fifth_book.year << "\n";


}