# A Go library to compute the contour of simple rectilinear polygons

This is not intended to be a general purpose library, but fits the need I had, namely to find the contour of a bunch of rectangles of possibly different sizes and mostly touching each other.

The heart of the library is the implementation of one of the sweep algorithms found in the following article :

> Souvaine, Diane & Bjorling-Sachs, Iliana. (1992). The Contour Problem for Restricted-Orientation Polygons. Proceedings of the IEEE. 80. 1449 - 1470. 10.1109/5.163411.
