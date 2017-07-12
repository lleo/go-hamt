Note: I have upgraded this package to add in the a functional mode as well.
This is a merger with github.com/lleo/go-hamt-functional which it obsoletes.
The motivation for this was because I was seeing slower performance for 
go-hamt-functional Get() operations even though I was certain that they were
using the same algorithm. This merger guarantees that the transient and
functinal Hamt implementations are useing the exact same datastructures. This
is true even to the degree that I can recast a pointer to HamtTransient or
HamtFunctional data structure and the code will switch from transient to
functional behavior (or vice versa).

It also obsoletes github.com/lleo/go-hamt-key because it need not be a separate
repository.

# What is a HAMT?

HAMT stands for Hash Array Mapped Trie. That spells it out clearly right? Ok,
not so much. A HAMT is a Key/Value storage with some really cool features. The
first feature is that we implement it as a tree with a reasonably largs and
fixed branching factor; thereby making it a pretty flat tree. The killer feature
is that the tree has a maximum depth, resulting in a O(1) Search and Modify
speeds. That is O(1) with out a giant constant factor.

HAMT make this happen by first hashing the key into either a 32bit (hamt32) or
64bit (hamt64) hash value. We then split this hash value up into a fixed number
of parts. Each part now becomes the index into an Array we use as interior nodes
of the tree. Lastly, we only use enough parts of the hash value to find a free
node to store the value into; hence Trie in its name. So now we have a wide tree
with a fixed depth where we only use enough of the hash value to find a unique
spot to store the value; That is what makes it O(1) and fast.

Specifically, we use the FNV hashing algorithm. It is fast and provides good
randomness and that is all we need from it.

Also we choose to split the 32bit or 64bit hash value into 5bit values. 5bit
values means the tree will have a branching factor of 32 and a fixed depth
of 6 for a 32bit hash value (hamt32) and 10 for a 64bit hash value (hamt64).
You made be noticing that 5 does not go into 32 nor 64 cleanly. That is not
a problem we fold the extra 2 or 4 bits into the main hash value. Don't worry
this is a legitimate thing to do. In the (very) rare case of a hash collision
we use a special leaf value for both colliding key/value pairs.

# go-hamt

We implement HAMT data structure based on either a 32 bit or 64 bit hash value,
hamt32 and hamt64 respectively.

Further we can have the HAMT data structure behave in one of two modes Transient
or Functional. Tansient means we modify the datastructures in-place. Functional
means use a copy-on-write strategy (ie. we copy each datastructure and modify
the copy). As you can imagine Tansient is faster than Functional. However, you
cannot easily share Transient datastructures between threads safely; you would
need to implement a locking strategy. Where with Functional data structures you
are guaranteed safty across threads.

## Functional (aka Immutable & Persistent)

The term Functional really stands for two properties Immutable and Persistent.
Immutable means the datastructure is never modified after construction.
Persistent means that when you modify the original Immutable datastructure you
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
