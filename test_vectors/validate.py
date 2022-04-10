#!/usr/bin/env python3

import json
import binascii
from os import path

# pip3 install ECPy
from ecpy.curves import Curve, Point
from ecpy.keys import ECPrivateKey
from ecpy.ecdsa import ECDSA

curve = Curve.get_curve('secp256k1')

def load_vectors(fname):
  with open(path.join(path.dirname(__file__), fname)) as fh:
    return json.load(fh)


for vector in load_vectors('negated_point.json'):
  expected_even_y = int(vector['evenY'], 16)
  expected_odd_y = int(vector['oddY'], 16)
  actual_odd_y = Point(int(vector['x'], 16), expected_even_y, curve).neg().y
  actual_even_y = Point(int(vector['x'], 16), expected_odd_y, curve).neg().y
  assert actual_even_y == expected_even_y
  assert actual_odd_y == expected_odd_y


for vector in load_vectors('jacobi_point.json'):
  expected_x = int(vector['x'], 16)
  expected_y = int(vector['y'], 16)
  jacobi_x = int(vector['jacobiX'], 16)
  jacobi_y = int(vector['jacobiY'], 16)
  jacobi_z = int(vector['jacobiZ'], 16)
  x, y = curve._jac2aff(jacobi_x, jacobi_y, jacobi_z, curve.field)
  assert x == expected_x
  assert y == expected_y


for vector in load_vectors('jacobi_multiplication.json'):
  expected_x = int(vector["x2"], 16)
  expected_y = int(vector["y2"], 16)
  p = Point(int(vector['x1'], 16), int(vector['y1'], 16), curve).mul(int(vector['k'], 16))
  if expected_x == 0 and expected_y == 0:
    assert p.is_infinity
  else:
    assert p.x == expected_x
    assert p.y == expected_y


for vector in load_vectors('jacobi_doubling.json'):
  jacobi_x = int(vector['x1'], 16)
  jacobi_y = int(vector['y1'], 16)
  jacobi_z = int(vector['z1'], 16)
  expected_x = int(vector['x3'], 16)
  expected_y = int(vector['y3'], 16)
  expected_z = int(vector['z3'], 16)
  x, y, z = curve._dbl_jac(jacobi_x, jacobi_y, jacobi_z, curve.field, curve.a)
  assert x == expected_x
  assert y == expected_y
  assert z == expected_z


# have to define this manually because ecpy doesn't expose safe jacobi addition directly
def add_jacobi(x1, y1, z1, x2, y2, z2):
  x1_aff, y1_aff = curve._jac2aff(x1, y1, z1, curve.field)
  x2_aff, y2_aff = curve._jac2aff(x2, y2, z2, curve.field)
  if x1_aff == 0 and y1_aff == 0: # 0 + Q = Q
    return x2, y2, z2
  elif x2_aff == 0 and y2_aff == 0: # P + 0 = P
    return x1, y1, z1
  elif x1_aff == x2_aff and y1_aff == y2_aff:
    return curve._dbl_jac(x1, y1, z1, curve.field, curve.a)
  return curve._add_jac(x1, y1, z1, x2, y2, z2, curve.field)

for vector in load_vectors('jacobi_addition.json'):
  x1 = int(vector['x1'], 16)
  y1 = int(vector['y1'], 16)
  z1 = int(vector['z1'], 16)
  x2 = int(vector['x2'], 16)
  y2 = int(vector['y2'], 16)
  z2 = int(vector['z2'], 16)
  expected_x3 = int(vector['x3'], 16)
  expected_y3 = int(vector['y3'], 16)
  expected_z3 = int(vector['z3'], 16)
  x3, y3, z3 = add_jacobi(x1, y1, z1, x2, y2, z2)

  # must compare affine values because ecpy uses slightly different
  # formulas resulting in different jacobian ratios
  expected_x3_aff, expected_y3_aff = curve._jac2aff(expected_x3, expected_y3, expected_z3, curve.field)
  x3_aff, y3_aff = curve._jac2aff(x3, y3, z3, curve.field)
  assert x3_aff == expected_x3_aff
  assert y3_aff == expected_y3_aff


ecdsa = ECDSA('ITUPLE')
for vector in load_vectors('ecdsa.json'):
  key = ECPrivateKey(int(vector["d"], 16), curve)
  h = binascii.unhexlify(vector['hash'])
  nonce = int(vector['nonce'], 16)
  expected_r = int(vector['r'], 16)
  expected_s = int(vector['s'], 16)
  r, s = ecdsa.sign_k(h, key, nonce, canonical=True)
  assert r == expected_r
  assert s == expected_s

print('All test vectors successfully validated!')
