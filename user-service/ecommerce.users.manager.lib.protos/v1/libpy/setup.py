from io import open
from os import path, pardir
from setuptools import setup, find_packages

parent_dir = path.abspath(path.join(path.dirname(__file__), pardir))

try:
    with open(path.join(parent_dir, "VERSION"), encoding="utf-8") as version_file:
        version = version_file.read().strip()
except FileNotFoundError:
    raise Exception("The VERSION file was not found")

with open(path.join(parent_dir, "README.md"), encoding="utf-8") as f:
    long_description = f.read()

setup(
    name="users_manager",
    version=version,
    author="Wembie",
    description="A package to manage users service with gRPC and protobuf",
    long_description=long_description,
    long_description_content_type="text/markdown",
    license="MIT",
    packages=find_packages(where="v1/libpy"),
    package_dir={"": "v1/libpy"},
    include_package_data=True,
    package_data={"users_manager": ["*.pyi"]},
    install_requires=[
        "grpcio>=1.71.0",
        "protobuf>=6.30.2",
    ],
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Developers",
        "Topic :: Software Development :: Libraries",
    ],
    python_requires=">=3.8",
    #url="https://github.com/Wembie/ecommerce-microservices",
)
