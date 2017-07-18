/*
Package hamt64 defines interface to access a Hamt data structure based on
64bit hash values. The Hamt data structure is built with interior nodes and leaf
nodes. The interior nodes are called tables and the leaf nodes are call, well,
leafs. Further, the tables come is two varieties fixed size tables and a
compressed form to handle sparse tables. Leafs come in two forms the common flat
leaf form with a singe key/value pair and the rare form used when two leafs have
the same hash value called collision leafs.

The Hamt data structure is implemented with two code bases, which both implement
the hamt64.Hamt interface, the transient replace in place code and the
functional copy on write code. We define a HamtTransient base data structure and
a HamtFunctional base data structure. Both of these data structures are
identical, they only have unique names so we can hang the different code
implementations off them.

Lastly, the Hamt data structure can be implemented with fixed tables only or
with sparse tables only or with a hybrid of the two. Thia hybid form is meant
to allow the denser lower inner nodes to be implemented by the faster fixed
tables and the much more numerous but sparser higher inner nodes to be
implemented by the space conscious sparse tables.
*/
package hamt64
