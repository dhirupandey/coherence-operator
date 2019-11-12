/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package com.oracle.coherence.k8s;

import java.net.UnknownHostException;
import java.util.Arrays;
import java.util.Iterator;

import com.tangosol.net.ConfigurableAddressProvider.AddressHolder;
import com.tangosol.util.Base;

import com.oracle.common.util.Duration;
import org.junit.Test;

import static com.oracle.coherence.k8s.RetryingWkaAddressProvider.PROP_WKA_OVERRIDE;
import static com.oracle.coherence.k8s.RetryingWkaAddressProvider.PROP_WKA_RERESOLVE_FREQUENCY;
import static com.oracle.coherence.k8s.RetryingWkaAddressProvider.PROP_WKA_TIMEOUT;
import static com.oracle.coherence.k8s.RetryingWkaAddressProvider.create;
import static com.oracle.common.util.Duration.Magnitude.MILLI;
import static org.hamcrest.Matchers.greaterThanOrEqualTo;
import static org.hamcrest.Matchers.is;
import static org.hamcrest.Matchers.lessThanOrEqualTo;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertThat;
import static org.junit.Assert.fail;

public class RetryingWkaAddressProviderTest {
    @Test
    public void testShouldTimeoutOnNonExistentDnsReference() {
        final long frequencyMS = 500;
        final long timeoutMS = 5000;

        long start = Base.getLastSafeTimeMillis();

        Iterable<AddressHolder> holders = new Iterable<AddressHolder>() {
            @Override
            public Iterator<AddressHolder> iterator() {
                return Arrays.asList(new AddressHolder("NonExiStentHoStName12345678", 0)).iterator();
            }

        };
        RetryingWkaAddressProvider provider = new RetryingWkaAddressProvider(holders, true, frequencyMS, timeoutMS);
        try {
            provider.eventuallyResolve();
            fail("should throw exception");
        }
        catch (UnknownHostException e) {
            // expected result.
        }
        finally {
            long timeoutDuration = Base.getLastSafeTimeMillis() - start;

            assertThat("validate reresolve happened over specified minimum timeout",
                       timeoutDuration, greaterThanOrEqualTo(timeoutMS));
            assertThat("validate lower bound of reresolve count",
                       provider.getLastReresolveCount(), greaterThanOrEqualTo(1));
            assertThat("validate upper bound of reresolve count",
                       provider.getLastReresolveCount(), lessThanOrEqualTo((int) (timeoutMS / frequencyMS)));
        }
    }

    @Test
    public void testShouldRevolveImmediately()
            throws UnknownHostException {
        final long frequencyMS = 3022;
        final long timeoutMS = 13456;

        RetryingWkaAddressProvider provider
                = (RetryingWkaAddressProvider) create("127.0.0.1", frequencyMS, timeoutMS);

        assertNotNull("confirm wka resolved", provider.getNextAddress());
        assertThat("validate configured frequency of reresolve", provider.getWkaDNSReresolveFrequency(), is(frequencyMS));
        assertThat("validate configured frequency of reresolve", provider.getWkaDNSResolutionTimeout(), is(timeoutMS));
    }

    @Test
    public void shouldConfigureBySystemProperties()
            throws UnknownHostException {
        String timeout = "4m";
        String frequency = "22s";

        System.setProperty(PROP_WKA_OVERRIDE, "127.0.0.1");
        System.setProperty(PROP_WKA_TIMEOUT, timeout);
        System.setProperty(PROP_WKA_RERESOLVE_FREQUENCY, frequency);

        try {
            RetryingWkaAddressProvider provider = (RetryingWkaAddressProvider) create();
            assertNotNull("confirm wka resolved", provider.getNextAddress());
            assertThat("validate configured frequency of dns resolve",
                       provider.getWkaDNSReresolveFrequency(), is(new Duration(frequency).as(MILLI)));
            assertThat("validate configured max time to attempt to resolve wka dns addresses",
                       provider.getWkaDNSResolutionTimeout(), is(new Duration(timeout).as(MILLI)));
        }
        finally {
            System.clearProperty(PROP_WKA_OVERRIDE);
            System.clearProperty(PROP_WKA_TIMEOUT);
            System.clearProperty(PROP_WKA_RERESOLVE_FREQUENCY);
        }
    }

    @Test
    public void testCreateWithDurationParameters()
            throws UnknownHostException {
        final String frequency = "2s";
        final String duration = "6s";

        System.setProperty(PROP_WKA_OVERRIDE, "127.0.0.1");

        try {
            RetryingWkaAddressProvider provider =
                    (RetryingWkaAddressProvider) create(frequency, duration);

            assertNotNull("confirm wka resolved", provider.getNextAddress());
            assertThat("validate configured frequency of reresolve",
                       provider.getWkaDNSReresolveFrequency(), is(new Duration(frequency).as(MILLI)));
            assertThat("validate configured frequency of reresolve",
                       provider.getWkaDNSResolutionTimeout(), is(new Duration(duration).as(MILLI)));
        }
        finally {
            System.clearProperty(PROP_WKA_OVERRIDE);
        }
    }
}