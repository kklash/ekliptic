use hex;
use libsecp256k1;
use serde_json;
use std::vec::Vec;
use std::{fs, io};

fn parse_u256(hex_str: &str) -> libsecp256k1::curve::Field {
    let mut bytes: [u8; 32] = [0; 32];
    match hex::decode_to_slice(hex_str, &mut bytes[..]) {
        Ok(()) => {}
        Err(msg) => {
            panic!("failed to parse {} as hex: {:?}", hex_str, msg);
        }
    };
    let mut n = libsecp256k1::curve::Field::new(0, 0, 0, 0, 0, 0, 0, 0);
    if !n.set_b32(&mut bytes) {
        panic!("failed to parse {} as hex", hex_str);
    };
    n
}

fn parse_scalar(hex_str: &str) -> libsecp256k1::curve::Scalar {
    let mut bytes: [u8; 32] = [0; 32];
    match hex::decode_to_slice(hex_str, &mut bytes[..]) {
        Ok(()) => {}
        Err(msg) => {
            panic!("failed to parse {} as hex: {:?}", hex_str, msg);
        }
    };
    let mut scalar = libsecp256k1::curve::Scalar([0, 0, 0, 0, 0, 0, 0, 0]);
    _ = scalar.set_b32(&mut bytes);
    scalar
}

fn to_affine(jacobian: &libsecp256k1::curve::Jacobian) -> libsecp256k1::curve::Affine {
    let mut affine = libsecp256k1::curve::Affine::from_gej(&jacobian);
    affine.x.normalize();
    affine.y.normalize();
    affine
}

fn read_vector_objects(fname: &str) -> Vec<serde_json::Value> {
    let file = fs::File::open("test_vectors/".to_owned() + fname).unwrap();
    let reader = io::BufReader::new(file);
    let vectors: serde_json::Value = serde_json::from_reader(reader).unwrap();
    let vectors_vec = vectors.as_array().unwrap();
    vectors_vec.to_vec()
}

// validates jacobi_point.json
fn validate_jacobi_points() {
    let vectors_vec = read_vector_objects("jacobi_point.json");

    for v in vectors_vec.iter() {
        let vector = v.as_object().unwrap();
        let jacobian = libsecp256k1::curve::Jacobian {
            x: parse_u256(&vector["jacobiX"].as_str().unwrap()),
            y: parse_u256(&vector["jacobiY"].as_str().unwrap()),
            z: parse_u256(&vector["jacobiZ"].as_str().unwrap()),
            infinity: false,
        };

        let expected_affine_x = vector["x"].as_str().unwrap();
        let expected_affine_y = vector["y"].as_str().unwrap();

        let affine = to_affine(&jacobian);
        let actual_affine_x = hex::encode(affine.x.b32());
        let actual_affine_y = hex::encode(affine.y.b32());

        if actual_affine_x != expected_affine_x || actual_affine_y != expected_affine_y {
            panic!(
                "jacobian point failed to convert to expected affine point:\nx: {}\ny: {}\nz: {}",
                hex::encode(jacobian.x.b32()),
                hex::encode(jacobian.y.b32()),
                hex::encode(jacobian.z.b32()),
            );
        }
    }
}

// validates negated_point.json
fn validate_negated_points() {
    let vectors_vec = read_vector_objects("negated_point.json");
    for v in vectors_vec.iter() {
        let vector = v.as_object().unwrap();
        let mut affine_even_y = libsecp256k1::curve::Affine {
            x: parse_u256(&vector["x"].as_str().unwrap()),
            y: parse_u256(&vector["evenY"].as_str().unwrap()),
            infinity: false,
        };
        let mut affine_odd_y = affine_even_y.neg();
        affine_odd_y.y.normalize();

        let actual_odd_y = hex::encode(affine_odd_y.y.b32());
        let expected_odd_y = vector["oddY"].as_str().unwrap();
        if actual_odd_y != expected_odd_y {
            panic!(
                "negated point oddY coordinate is incorrect.\nWanted {}\nGot    {}",
                expected_odd_y, actual_odd_y,
            );
        }

        // negate again, back to evenY, and ensure nothing changed
        affine_even_y = affine_odd_y.neg();
        affine_even_y.y.normalize();

        let actual_even_y = hex::encode(affine_even_y.y.b32());
        let expected_even_y = vector["evenY"].as_str().unwrap();
        if actual_even_y != expected_even_y {
            panic!(
                "negated point evenY coordinate is incorrect.\nWanted {}\nGot    {}",
                expected_even_y, actual_even_y,
            );
        }
    }
}

// validates jacobi_multiplication.json
fn validate_jacobi_multiplication() {
    let ecmult_ctx = libsecp256k1::curve::ECMultContext::new_boxed();

    let vectors_vec = read_vector_objects("jacobi_multiplication.json");
    for v in vectors_vec.iter() {
        let vector = v.as_object().unwrap();
        let p1 = libsecp256k1::curve::Affine {
            x: parse_u256(&vector["x1"].as_str().unwrap()),
            y: parse_u256(&vector["y1"].as_str().unwrap()),
            infinity: false,
        };

        let scalar = parse_scalar(&vector["k"].as_str().unwrap());
        let mut result = libsecp256k1::curve::Jacobian::default();
        ecmult_ctx.ecmult_const(&mut result, &p1, &scalar);

        let mut affine_result = to_affine(&result);
        affine_result.x.normalize();
        affine_result.y.normalize();

        let actual_x2 = hex::encode(affine_result.x.b32());
        let actual_y2 = hex::encode(affine_result.y.b32());

        let expected_x2 = vector["x2"].as_str().unwrap();
        let expected_y2 = vector["y2"].as_str().unwrap();

        if actual_x2 != expected_x2 || actual_y2 != expected_y2 {
            panic!(
                "jacobi multiplication failed:\nx: {}\ny: {}\nk: {}",
                vector["x1"].as_str().unwrap(),
                vector["y1"].as_str().unwrap(),
                vector["k"].as_str().unwrap(),
            );
        }
    }
}

// validates jacobi_doubling.json
fn validate_jacobi_doubling() {
    let vectors_vec = read_vector_objects("jacobi_doubling.json");
    for v in vectors_vec.iter() {
        let vector = v.as_object().unwrap();
        let p1 = libsecp256k1::curve::Jacobian {
            x: parse_u256(&vector["x1"].as_str().unwrap()),
            y: parse_u256(&vector["y1"].as_str().unwrap()),
            z: parse_u256(&vector["z1"].as_str().unwrap()),
            infinity: false,
        };
        let mut p3 = p1.double_var(None);
        p3.x.normalize();
        p3.y.normalize();
        p3.z.normalize();

        let actual_x3 = hex::encode(p3.x.b32());
        let actual_y3 = hex::encode(p3.y.b32());
        let actual_z3 = hex::encode(p3.z.b32());

        let expected_x3 = vector["x3"].as_str().unwrap();
        let expected_y3 = vector["y3"].as_str().unwrap();
        let expected_z3 = vector["z3"].as_str().unwrap();

        if actual_x3 != expected_x3 || actual_y3 != expected_y3 || actual_z3 != expected_z3 {
            panic!(
                "failed to double to expected jacobi point:\nx: {}\ny: {}\nz: {}",
                vector["x1"].as_str().unwrap(),
                vector["y1"].as_str().unwrap(),
                vector["z1"].as_str().unwrap(),
            );
        }
    }
}

// validates jacobi_addition.json
fn validate_jacobi_addition() {
    let zero = parse_u256("0000000000000000000000000000000000000000000000000000000000000000");

    let vectors_vec = read_vector_objects("jacobi_addition.json");
    for v in vectors_vec.iter() {
        let vector = v.as_object().unwrap();

        let mut p1 = libsecp256k1::curve::Jacobian {
            x: parse_u256(&vector["x1"].as_str().unwrap()),
            y: parse_u256(&vector["y1"].as_str().unwrap()),
            z: parse_u256(&vector["z1"].as_str().unwrap()),
            infinity: false,
        };
        if p1.x.eq_var(&zero) && p1.y.eq_var(&zero) {
            p1.infinity = true;
        }

        let mut p2 = libsecp256k1::curve::Jacobian {
            x: parse_u256(&vector["x2"].as_str().unwrap()),
            y: parse_u256(&vector["y2"].as_str().unwrap()),
            z: parse_u256(&vector["z2"].as_str().unwrap()),
            infinity: false,
        };
        if p2.x.eq_var(&zero) && p2.y.eq_var(&zero) {
            p2.infinity = true;
        }

        let mut p3 = p1.add_var(&p2, None);
        p3.x.normalize();
        p3.y.normalize();
        p3.z.normalize();

        let actual_x3 = hex::encode(p3.x.b32());
        let actual_y3 = hex::encode(p3.y.b32());
        let actual_z3 = hex::encode(p3.z.b32());

        let expected_x3 = vector["x3"].as_str().unwrap();
        let expected_y3 = vector["y3"].as_str().unwrap();
        let expected_z3 = vector["z3"].as_str().unwrap();

        if actual_x3 != expected_x3 || actual_y3 != expected_y3 || actual_z3 != expected_z3 {
            panic!(
                "failed to add to expected jacobi point:\nx: {}\ny: {}\nz: {}\n  +\nx: {}\ny: {}\nz: {}",
                vector["x1"].as_str().unwrap(),
                vector["y1"].as_str().unwrap(),
                vector["z1"].as_str().unwrap(),
                vector["x2"].as_str().unwrap(),
                vector["y2"].as_str().unwrap(),
                vector["z2"].as_str().unwrap(),
            );
        }
    }
}

// validates ecdsa.json
fn validate_ecdsa() {
    let ecmult_ctx = libsecp256k1::curve::ECMultGenContext::new_boxed();

    let vectors_vec = read_vector_objects("ecdsa.json");
    for v in vectors_vec.iter() {
        let vector = v.as_object().unwrap();
        let d = parse_scalar(&vector["d"].as_str().unwrap());
        let hash = parse_scalar(&vector["hash"].as_str().unwrap());
        let nonce = parse_scalar(&vector["nonce"].as_str().unwrap());

        let (r, s, _) = ecmult_ctx.sign_raw(&d, &hash, &nonce).unwrap();
        let actual_r = hex::encode(r.b32());
        let actual_s = hex::encode(s.b32());

        let expected_r = vector["r"].as_str().unwrap();
        let expected_s = vector["s"].as_str().unwrap();

        if actual_r != expected_r || actual_s != expected_s {
            panic!(
                "failed to validate ECDSA signature:\nr: {}\ns: {}\n\nGot:\nr: {}\ns: {}",
                expected_r, expected_s, actual_r, actual_s,
            );
        }
    }
}

fn main() {
    validate_jacobi_points();
    validate_negated_points();
    validate_jacobi_multiplication();
    validate_jacobi_doubling();
    validate_jacobi_addition();
    validate_ecdsa();

    println!("All test vectors successfully validated!");
}
