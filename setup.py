__author__ = "Kristoffer Dalby <kradalby@kradalby.no>"
__version__ = "0.0.0"

from setuptools import setup
import sys


setup(
    name="kubespace",
    version=__version__,
    description="Tool (kubectl wrapper) to create namespaces and service accounts that can safely be handed to CI or users.",
    author=__author__,
    author_email="kradalby@kradalby.no",
    url="https://github.com/kradalby/kubespace",
    license="MIT",
    # Ansible will also make use of a system copy of python-six and
    # python-selectors2 if installed but use a Bundled copy if it's not.
    # install_requires=install_requirements,
    scripts=["kubespace"],
    data_files=[],
    # Installing as zip files would break due to references to __file__
    zip_safe=False,
)
