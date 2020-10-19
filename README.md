# V2 API

This package was transient only, now it has a functional mode as well. Transient
is the classical style of modifying data structures in place. Functional is
defined below as immutable & persistent.

This is a merger with github.com/lleo/go-hamt-functional, which it obsoletes.
The motivation for this was because we were seeing slower performance for
go-hamt-functional Get() operations even though we were certain that the old
go-hamt and go-hamt-functional were using the same algorithm. This merger
guarantees that the transient and functional Hamt implementations are using the
exact same internal data structures. This is true even to the degree that we can
recast a HamtTransient data structure to HamtFunctional and the code will switch
from transient (modify in place) to functional (copy on write) behavior. Of
course, this works the other way around as well (that is, we can cast a
HamtFunctional to HamtTransient).

This package also obsoletes github.com/lleo/go-hamt-key because we pass a []byte
slice to Get/Put/Del operations instead of a Key data structure. What happens
is we use the []byte slice to build a Key data structure to be used internally.
This results in a simpler API and no external dependency.

## What is a HAMT?

HAMT stands for Hash Array Mapped Trie. That spells it out clearly right? Ok,
not so much. A HAMT is an in-memory Key/Value data structure with some really
cool features. The first feature is that we implement it as a tree with a
reasonably large and fixed branching factor; thereby making it a pretty flat
tree. The killer feature is that the tree has a maximum depth, resulting in a
O(1) Search and Modify speeds. That is O(1) without a giant constant factor.

HAMT make this happen by first hashing the []byte slice  into either a 32bit or
64bit hash value. We then split this hash value up into a fixed number of parts.
Each part now becomes the index into the Array we use for interior nodes of the
tree. Lastly, we only use enough parts of the hash value to find a free node to
store the value into; this is why it is called a Trie. So now we have a wide
tree with a maximum depth where we only use enough of the parts of hash value to
find a free spot to store the leaf; That is what makes a HAMT O(1) and fast.

In our implementation, we use the FNV hashing algorithm. It is fast and provides
good randomness and that is all we need from it.

Also we choose to split the 32bit or 64bit hash value into 5bit values. 5bit
values mean the tree will have a branching factor of 32 and a maximum depth of
6 for a 32bit hash value (hamt32) and 12 for a 64bit hash value (hamt64). You
may be noticing that 5 does not go into 32 nor 64 cleanly. That is not a problem
because we fold the extra 2 or 4 bits into the main hash value. Don't worry this
is a legitimate thing to do. In the (very) rare case of a hash collision we use
a special leaf value for both colliding key/value pairs.

## go-hamt

We implement HAMT data structure based on either a 32 bit or 64 bit hash value,
hamt32 and hamt64 respectively.

Further we can have the HAMT data structure behave in one of two modes,
transient or functional. Transient means we modify the data structures in-place.
Functional means persistent, which requires we use a copy-on-write strategy. In
other words we copy each data structure and modify the copy, then given the
parent now needs to be modified we follow this copy-on-write strategy up to a
brand new HamtFunctional data structure. As you can imagine Transient is faster
than Functional. Suprisingly, the transient behavior is not that much faster
than the persistent behavior, because the HAMT data structure is flat and wide.

However, you cannot easily share transient datastructures between threads
safely; you would need to implement a locking strategy. Where with functional
data structures you are guaranteed safety across threads.

### Functional (aka Immutable & Persistent)

The term Functional really stands for two properties immutable and persistent.
Immutable means the datastructure is never modified after construction.
Persistent means that when you modify the original immutable datastructure you
do so by creating a new copy of the base datastructure which shares all its
unmodified parts with the original datastructure.

Imagine a hypothetical balanced binary tree datastructure with four leaves, two
interior nodes, and a root node. If you change the forth leaf node, than a new
fourth tree node is created, as well as its parent interior node and a new root
node. The new tree starting at the new root node is a persistent modification of
the original tree as it shares the first interior node and its two leaves.

              (root tree node)     (root tree node')
                /          \         /          \
               /  +---------\-------+            \
              /  /           \                    \
        (tree node 1)    (tree node 2)        (tree node 2')
            /  \            /  \                /   \
           /    \          / +--\--------------+     \
          /      \        / /    \                    \
     (Leaf 1) (Leaf 2) (Leaf 3) (Leaf 4)             (Leaf 4')

Given this approach to changing a tree, a tree with a wide branching factor
would be relatively shallow. So the path from root to any given leaf would be
short and the amount of shared content would be substantial.
