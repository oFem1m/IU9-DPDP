cmake_minimum_required(VERSION 3.29)
project(IU9_DPDP_GO)

set(CMAKE_CXX_STANDARD 14)

find_package(OpenMP REQUIRED)
if(OpenMP_CXX_FOUND)
    add_executable(lab3 lab3/lab3.cpp)
    target_link_libraries(lab3 PUBLIC OpenMP::OpenMP_CXX)
else()
    message(FATAL_ERROR "OpenMP not found")
endif()
